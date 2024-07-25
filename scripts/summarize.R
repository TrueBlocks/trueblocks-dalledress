library(ggplot2)
library(dplyr)

print(getwd())
setwd("output/selectors/")

files <- list.files(pattern = "\\.txt$")
data_list <- lapply(files, read.table)
names(data_list) <- files

data_combined <- do.call(rbind, lapply(names(data_list), function(name) {
  df <- data_list[[name]]
  df$attribute <- name
  df
}))

ggplot(data_combined, aes(x = V1)) +
  geom_histogram(binwidth = 1, fill = "blue", color = NA) +
  facet_wrap(~ attribute, scales = "free") +
  labs(title = "Distribution of Values Across Attributes", x = "Value", y = "Frequency") +
  theme_minimal()

ggplot(data_combined, aes(x = V1, fill = attribute)) +
  geom_density(alpha = 0.5, adjust = 1, trim = TRUE) +
  facet_wrap(~ attribute, scales = "free") +
  labs(title = "Density of Values Across Attributes", x = "Value", y = "Density") +
  theme_minimal()


summary_stats <- data_combined %>%
  group_by(attribute) %>%
  summarize(
    mean = mean(V1),
    median = median(V1),
    sd = sd(V1),
    min = min(V1),
    max = max(V1)
  )

print(summary_stats)

summary_stats %>%
  mutate(diff_mean_median = abs(mean - median)) %>%
  select(attribute, mean, median, diff_mean_median)

summary_stats %>%
  mutate(cv = sd / mean) %>%
  select(attribute, mean, sd, cv)

ggplot(data_combined, aes(x = V1)) +
  geom_histogram(binwidth = 10, fill = "blue", color = "black") +
  facet_wrap(~ attribute, scales = "free") +
  labs(title = "Distribution of Values Across Attributes", x = "Value", y = "Frequency") +
  theme_minimal()

ggplot(data_combined, aes(y = V1, x = attribute)) +
  geom_boxplot() +
  labs(title = "Boxplots of Values Across Attributes", x = "Attribute", y = "Value") +
  theme_minimal() +
  theme(axis.text.x = element_text(angle = 90, hjust = 1))


count_data <- data_combined %>%
  group_by(attribute, V1) %>%
  summarize(count = n(), .groups = 'drop')

stats_data <- count_data %>%
  group_by(attribute) %>%
  summarize(mean_count = mean(count), sd_count = sd(count))

count_data <- merge(count_data, stats_data, by = "attribute")

summary_stats <- count_data %>%
  group_by(attribute) %>%
  summarize(
    mean_count = mean(count),
    sd_count = sd(count)
  )
summary_stats

ggplot(count_data, aes(x = count)) +
  geom_histogram(aes(y = ..density..), binwidth = 1, fill = "blue", color = "black") +
  stat_function(aes(color = attribute), fun = function(x) dnorm(x, mean = mean_count, sd = sd_count), linetype = "dashed", size = 0.5) +
  facet_wrap(~ attribute, scales = "free") +
  labs(title = "Distribution of Value Counts Across Attributes with Normal Curve", x = "Count of Values", y = "Density") +
  theme_minimal()






summary_stats <- data.frame(
  attribute = c("action.txt", "adjective.txt", "adverb.txt", "artStyle1.txt", "artStyle2.txt", "backStyle.txt", 
                "color1.txt", "color2.txt", "color3.txt", "emotion.txt", "gaze.txt", "litStyle.txt", 
                "noun.txt", "occupation.txt", "orientation.txt", "setting.txt"),
  mean_count = c(13.2, 10.6, 4.55, 42.1, 42.1, 1711, 98.5, 98.5, 98.5, 49.1, 1711, 144, 4.02, 62.2, 1711, 7.61),
  sd_count = c(3.63, 3.35, 2.14, 6.74, 6.65, 49.7, 10, 9.93, 9.13, 7.87, 45.9, 13.5, 1.97, 7.4, 42.8, 2.86)
)

count_data <- count_data %>%
  left_join(summary_stats, by = "attribute")

ggplot(count_data, aes(x = count)) +
  geom_histogram(aes(y = ..density..), binwidth = 1, fill = "blue", color = "black") +
  stat_function(
    aes(color = attribute),
    fun = function(x, mean_count, sd_count) dnorm(x, mean = mean_count, sd = sd_count),
    args = list(mean_count = summary_stats$mean_count, sd_count = summary_stats$sd_count),
    linetype = "dashed", size = 0.5
  ) +
  facet_wrap(~ attribute, scales = "free") +
  labs(title = "Distribution of Value Counts Across Attributes with Normal Curve", x = "Count of Values", y = "Density") +
  theme_minimal() +
  theme(axis.text.x = element_text(angle = 90, hjust = 1))





count_data <- data.frame(
  attribute = c("action.txt", "adjective.txt", "adverb.txt", "artStyle1.txt", "artStyle2.txt", "backStyle.txt", 
                "color1.txt", "color2.txt", "color3.txt", "emotion.txt", "gaze.txt", "litStyle.txt", 
                "noun.txt", "occupation.txt", "orientation.txt", "setting.txt"),
  count = c(13.2, 10.6, 4.55, 42.1, 42.1, 1711, 98.5, 98.5, 98.5, 49.1, 1711, 144, 4.02, 62.2, 1711, 7.61)
)

count_data <- count_data %>%
  mutate(log_count = log(count))

# If there are zeros or negative values in the data
count_data <- count_data %>%
  mutate(log_count = log1p(count))  # log1p(x) = log(1 + x)

# Plot the transformed data
ggplot(count_data, aes(x = log_count, fill = attribute)) +
  geom_density(alpha = 0.5) +
  facet_wrap(~ attribute, scales = "free") +
  labs(title = "Log-Transformed Density Plots of Counts Across Attributes", x = "Log Count", y = "Density") +
  theme_minimal()





library(dplyr)
library(ggplot2)

# Example data (assuming count_data is your original dataset)
# If your actual data is different, replace this with your real data loading step
count_data <- data.frame(
  attribute = rep(c("action.txt", "adjective.txt", "adverb.txt", "artStyle1.txt", "artStyle2.txt", "backStyle.txt", 
                    "color1.txt", "color2.txt", "color3.txt", "emotion.txt", "gaze.txt", "litStyle.txt", 
                    "noun.txt", "occupation.txt", "orientation.txt", "setting.txt"), each = 100),
  count = c(
    rnorm(100, 13.2, 3.63), rnorm(100, 10.6, 3.35), rnorm(100, 4.55, 2.14), 
    rnorm(100, 42.1, 6.74), rnorm(100, 42.1, 6.65), rnorm(100, 1711, 49.7), 
    rnorm(100, 98.5, 10.0), rnorm(100, 98.5, 9.93), rnorm(100, 98.5, 9.13), 
    rnorm(100, 49.1, 7.87), rnorm(100, 1711, 45.9), rnorm(100, 144, 13.5), 
    rnorm(100, 4.02, 1.97), rnorm(100, 62.2, 7.40), rnorm(100, 1711, 42.8), 
    rnorm(100, 7.61, 2.86))
)

# Handle negative values or zero by adding a small constant before log transformation
# count_data <- count_data %>% mutate(log_count = log(count + 1)) # Use log1p if you have zeros

# Apply log transformation
count_data <- count_data %>%
  mutate(log_count = ifelse(count <= 0, NA, log(count)))

# Plot the transformed data
ggplot(count_data, aes(x = log_count, fill = attribute)) +
  geom_density(alpha = 0.5) +
  facet_wrap(~ attribute, scales = "free") +
  labs(title = "Log-Transformed Density Plots of Counts Across Attributes", x = "Log Count", y = "Density") +
  theme_minimal()


